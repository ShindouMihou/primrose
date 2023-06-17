export const Dates = {
    months: ["January","February","March","April","May","June","July","August","September","October","November","December"],
    toDateString: (date: Date): string => {
        return Dates.months[date.getMonth()] + " " +  date.getDate() + ", " + date.getFullYear()
    },
}