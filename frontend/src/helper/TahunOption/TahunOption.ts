export const tahunOption = () => {
    const current = new Date().getFullYear();
    const start = 2022
    const options = []

    for (let i = current; i >= start; i--){
        options.push(i)
    }

    return options
}