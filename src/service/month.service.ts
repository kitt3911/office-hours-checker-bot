import { prisma } from "../config"
import { Day } from "../interfaces/day.interface"
import { formatDate } from "../utils/formatDate"
import { validationMonth } from "./validation.service"

export const getDayOnMonthService = async(userId?:string,monthId?: string):Promise<Day[]>=> {
    const days = await prisma.day.findMany({
        where: {
            monthId:monthId
        }
    })
    return days
}