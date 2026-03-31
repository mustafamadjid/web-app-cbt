package gradingrepo

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type execCall struct {
	query string
	args  []any
}

type execResult struct {
	tag pgconn.CommandTag
	err error
}

type fakeExecutor struct {
	execResults []execResult
	execCalls   []execCall
}

func (f *fakeExecutor) Exec(_ context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	f.execCalls = append(f.execCalls, execCall{query: sql, args: arguments})

	idx := len(f.execCalls) - 1
	if idx >= len(f.execResults) {
		return pgconn.NewCommandTag(""), errors.New("unexpected exec call")
	}

	result := f.execResults[idx]
	return result.tag, result.err
}

func (*fakeExecutor) Query(context.Context, string, ...any) (pgx.Rows, error) {
	return nil, errors.New("unexpected query call")
}

func (*fakeExecutor) QueryRow(context.Context, string, ...any) pgx.Row {
	return nil
}

func TestUpsertNilaiToHasilUjian_UsesSingleUpsertQuery(t *testing.T) {
	t.Parallel()

	executor := &fakeExecutor{
		execResults: []execResult{
			{tag: pgconn.NewCommandTag("INSERT 0 1")},
		},
	}

	repo := NewGradingRepo(executor, nil)

	err := repo.UpsertNilaiToHasilUjian(context.Background(), 87.5, ujian.HasilUjian{
		IdAttempt: 12,
	})

	require.NoError(t, err)
	require.Len(t, executor.execCalls, 1)
	assert.Contains(t, executor.execCalls[0].query, "INSERT INTO hasil_ujian")
	assert.Contains(t, executor.execCalls[0].query, "ON CONFLICT (id_attempt)")
	assert.Contains(t, executor.execCalls[0].query, "graded_by = COALESCE($2, hasil_ujian.graded_by)")
	assert.Contains(t, executor.execCalls[0].query, "passed = COALESCE($4, hasil_ujian.passed)")
	assert.Contains(t, executor.execCalls[0].query, "essay_graded = COALESCE($5, hasil_ujian.essay_graded)")
	assert.Contains(t, executor.execCalls[0].query, "graded_at = COALESCE($6, NOW())")

	require.Len(t, executor.execCalls[0].args, 6)
	assert.Equal(t, ujian.ID(12), executor.execCalls[0].args[0])
	assert.Nil(t, executor.execCalls[0].args[1])
	assert.Equal(t, 87.5, executor.execCalls[0].args[2])
	assert.Nil(t, executor.execCalls[0].args[3])
	assert.Nil(t, executor.execCalls[0].args[4])
	assert.Nil(t, executor.execCalls[0].args[5])
}

func TestUpsertNilaiToHasilUjian_PassesOptionalFields(t *testing.T) {
	t.Parallel()

	executor := &fakeExecutor{
		execResults: []execResult{
			{tag: pgconn.NewCommandTag("INSERT 0 1")},
		},
	}

	repo := NewGradingRepo(executor, nil)

	gradedBy := ujian.ID(7)
	passed := true
	essayGraded := true
	gradedAt := time.Date(2026, time.March, 31, 13, 33, 1, 0, time.UTC)

	err := repo.UpsertNilaiToHasilUjian(context.Background(), 91.25, ujian.HasilUjian{
		IdAttempt:   21,
		GradedBy:    &gradedBy,
		Passed:      &passed,
		EssayGraded: &essayGraded,
		GradedAt:    &gradedAt,
	})

	require.NoError(t, err)
	require.Len(t, executor.execCalls, 1)
	assert.Equal(t, []any{ujian.ID(21), &gradedBy, 91.25, &passed, &essayGraded, &gradedAt}, executor.execCalls[0].args)
}

func TestUpsertJawabanBenarToStatistikSoal_SQLAggregatesDuplicateRows(t *testing.T) {
	t.Parallel()

	executor := &fakeExecutor{
		execResults: []execResult{
			{tag: pgconn.NewCommandTag("INSERT 0 2")},
		},
	}

	repo := NewGradingRepo(executor, nil)

	err := repo.UpsertJawabanBenarToStatistikSoal(context.Background(), []ujian.StatistikSoal{
		{IDUjian: 3, IDSoal: 10},
		{IDUjian: 3, IDSoal: 10},
		{IDUjian: 3, IDSoal: 11},
	})

	require.NoError(t, err)
	require.Len(t, executor.execCalls, 1)
	assert.Contains(t, executor.execCalls[0].query, "WITH aggregated AS")
	assert.Contains(t, executor.execCalls[0].query, "SUM(x.jumlah_jawaban_benar)")
	assert.Contains(t, executor.execCalls[0].query, "GROUP BY x.id_soal, x.id_ujian")

	require.Len(t, executor.execCalls[0].args, 1)
	payloadJSON, ok := executor.execCalls[0].args[0].([]byte)
	require.True(t, ok)

	var payload []statistikSoalUpsertItem
	require.NoError(t, json.Unmarshal(payloadJSON, &payload))
	require.Len(t, payload, 3)

	assert.Equal(t, statistikSoalUpsertItem{
		IDSoal:             10,
		IDUjian:            3,
		JumlahJawabanBenar: 1,
	}, payload[0])
	assert.Equal(t, statistikSoalUpsertItem{
		IDSoal:             10,
		IDUjian:            3,
		JumlahJawabanBenar: 1,
	}, payload[1])
	assert.Equal(t, statistikSoalUpsertItem{
		IDSoal:             11,
		IDUjian:            3,
		JumlahJawabanBenar: 1,
	}, payload[2])
}

func TestUpsertJawabanSalahToStatistikSoal_EmptySliceNoop(t *testing.T) {
	t.Parallel()

	executor := &fakeExecutor{}
	repo := NewGradingRepo(executor, nil)

	err := repo.UpsertJawabanSalahToStatistikSoal(context.Background(), nil)

	require.NoError(t, err)
	assert.Empty(t, executor.execCalls)
}
