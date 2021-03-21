import { PrismaClient } from "@prisma/client";
import { Context } from "telegraf";
import { prisma } from "../config";
import { Day } from "../interfaces/day.interface";
import { Month } from "../interfaces/month.interface";
import { User } from "../interfaces/user.interface";
import { formatDate } from "../utils/formatDate";

export const validationUser = async (user: User): Promise<User> => {
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

export const validationMonth = async (userId?: string,date?: string):Promise<Month> => {
    let thisDate = formatDate(date)
    let findMonth = await prisma.month.findFirst({
        where: {
            fullDate: thisDate.monthAndYear,
            userId
        }
    })
    if (!findMonth) {
        findMonth = await prisma.month.create({
          data: {
              fullDate: thisDate.monthAndYear,
              userId,
              name: thisDate.month
          }
      })
  }

  return findMonth
}

