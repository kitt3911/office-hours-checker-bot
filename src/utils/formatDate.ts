import moment from 'moment'

export const formatMonth = (date?: string) => {
    let thisDate: moment.Moment
    if (!date) {
        thisDate = moment(new Date(), 'DD/MM/YYYY')
    }
    else {
        let dateForKnowYear = new Date()
        let buffDate = new Date(date + '/' + dateForKnowYear.getFullYear())
        thisDate = moment(buffDate)
    }
    let month = thisDate.format('M')
    let year = thisDate.format('YYYY')
    let thisFormatDate = year + '/' + month
    return {
        date: thisFormatDate, month, year
    }
}