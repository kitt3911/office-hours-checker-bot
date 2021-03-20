import { PrismaClient } from "@prisma/client"
import { Day } from "../interfaces/day.interface"
import { formatMonth } from "../utils/formatDate"

export const findOneDay = async (prisma: PrismaClient,monthId?: string):Promise<Day|null> => {
    let findDay = await prisma.day.findFirst({
        where: {
            monthId
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
    const formatDate = formatMonth(date)
    const formatDay = formatDate.date + '/' + formatDate.day
    if(findDay){
        return await prisma.day.update({
           where: {
                id:findDay.id
           },
           data: {
               date: formatDay,
               workHours
           }
        })
    }

    return await prisma.day.create({
        data: {
            monthId,
            date: formatDay,
            workHours
        }
    })
}