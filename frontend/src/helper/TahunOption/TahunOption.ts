export const tahunOption = () => {
    const current = new Date().getFullYear();
    const start = 2019
    const options = []

    for (let i = current; i >= start; i--){
        options.push(String(i))
    }

    return options
}