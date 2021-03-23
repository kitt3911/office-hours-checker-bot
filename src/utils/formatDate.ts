import moment from 'moment'
import { FormatDate } from '../interfaces/formatDate.interface'

export const formatDate = (date?: string):FormatDate => {
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
    let thisFormatDayDate = thisFormatDate + '/' + day
    return {
        monthAndYear: thisFormatDate,fullDate:thisFormatDayDate, month, year, day
    }
}

export const formatHours = (workHours: number): string => {
    const hours = Math.trunc(workHours)
    const minNumber = ((workHours - hours)*60)
    let min = minNumber.toString()
    if(min.length === 1){
        min = min + '0'
    }

    return `${hours}:${min}`
}

export const convertMinToDecimal = (str: string) : number => {
    return Number(str) / 60;
}