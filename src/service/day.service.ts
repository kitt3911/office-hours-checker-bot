import { PrismaClient } from "@prisma/client"
import { prisma } from "../config"
import { Day } from "../interfaces/day.interface"
import { formatDate } from "../utils/formatDate"

export const findOneDay = async (monthId?: string,date?: string):Promise<Day|null> => {
    const format = formatDate(date)
    let findDay = await prisma.day.findFirst({
        where: {
            monthId,
            date: format.fullDate
        }
    })
    return findDay
}

export const createDay = async (monthId: string,workHours: number,date: string): Promise<Day> => {

    const format = formatDate(date)
    const findDay = await prisma.day.findFirst({
        where: {
            monthId,
            date: format.fullDate
        }
    })
    console.log(findDay)
    if(findDay){
        return await prisma.day.update({
           where: {
                id:findDay.id
           },
           data: {
               date: format.fullDate,
               workHours
           }
        })
    }

    return await prisma.day.create({
        data: {
            monthId,
            date: format.fullDate,
            workHours
        }
    })
}