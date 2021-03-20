import moment from 'moment'
import { FormatDate } from '../interfaces/formatDate.interface'

export const formatMonth = (date?: string):FormatDate => {
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
    let day = thisDate.format('D')
    let thisFormatDate = year + '/' + month
    return {
        date: thisFormatDate, month, year, day
    }
}

export const formatHours = (workHours: number): string => {
    const min = Number(workHours.toFixed(0))/10
    const hours = Number.parseInt(workHours.toString())
    return `${hours}:${min}`
}