import { PrismaClient } from "@prisma/client";
import { Context } from "telegraf";
import { Day } from "../interfaces/day.interface";
import { Month } from "../interfaces/month.interface";
import { User } from "../interfaces/user.interface";
import { formatDate } from "../utils/formatDate";

export const validationUser = async (prisma: PrismaClient,user: User): Promise<User> => {
    let findUser = await prisma.user.findFirst({
        where: { telegramId: user.telegramId }
    })
    if (!findUser) {
        findUser = await prisma.user.create({
            data: { ...user }
        })
    }   
    return findUser
}

export const validationMonth = async (prisma: PrismaClient,date?: string,userId?: string):Promise<Month> => {
    let thisDate = formatDate(date)
    let findMonth = await prisma.month.findFirst({
        where: {
            fullDate: thisDate.localDate,
            userId
        }
    })
    if (!findMonth) {
        findMonth = await prisma.month.create({
          data: {
              fullDate: thisDate.localDate,
              userId,
              name: thisDate.month
          }
      })
  }

  return findMonth
}

