const join = (base: string, path: string) => `${base}${path.startsWith('/')? path : `/${path}`}`

export const paths = {
    public: {
        home: "/",
        login: "/login"
    }
} as const