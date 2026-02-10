package kelas_service

import (
	"context"
	"strings"

	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	txout "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/tx"
)

type CreateKelasTx struct {
	txm txout.TxManager
}

func NewCreateKelasService(txm txout.TxManager) *CreateKelasTx {
	return &CreateKelasTx{txm: txm}
}

func (tc *CreateKelasTx)CreateKelas(ctx context.Context, cmd CreateKelasCmd)error{
	 logger := corelog.FromContext(ctx)

	//  Normalisasi
	cmd.NamaKelas = strings.TrimSpace(cmd.NamaKelas)

	// Transaksi
	tx,err := tc.txm.Begin(ctx)
	if err != nil {
		logger.Error(ctx, "failed starting transaction", "layer", "core.service", "op", "kelas.create_kelas", "err", err)
	}

	defer func(){
		_ = tx.Rollback()
	}()
	return nil
}