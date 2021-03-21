import { PrismaClient } from "@prisma/client"
import { Day } from "../interfaces/day.interface"
import { formatDate } from "../utils/formatDate"

export const findOneDay = async (prisma: PrismaClient,monthId?: string,date?: string):Promise<Day|null> => {
    const format = formatDate(date)
    let findDay = await prisma.day.findFirst({
        where: {
            monthId,
            date: format.fullDate
        }
    })
    return findDay
}

export const createDay = async (prisma: PrismaClient,monthId: string,workHours: number,date: string): Promise<Day> => {

    const findDay = await prisma.day.findFirst({
        where: {
            monthId
        }
    })
    const format = formatDate(date)
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