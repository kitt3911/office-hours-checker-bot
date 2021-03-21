
import { PrismaClient } from "@prisma/client"
import { Telegraf } from "telegraf"
import { prisma } from "../config"
import { validationMonth, validationUser } from "../service/validation.service"

export const startController = async(bot: Telegraf)=>{
    bot.start(async (ctx) => {
        const user = {
            name: ctx.from.username,
            telegramId: ctx.from.id
        }
        const validateUser = await validationUser(user)
        const validateMonth = await validationMonth(validateUser.id)

        ctx.reply(`Welcome ${validateUser.name}`)
    })
}