import { PrismaClient } from "@prisma/client"
import { Telegraf } from "telegraf"
import { prisma } from "../config"
import { createDay, findOneDay } from "../service/day.service"
import { validationMonth, validationUser } from "../service/validation.service"
import { convertMinToDecimal, formatDate, formatHours } from "../utils/formatDate"

export const getDayController = async (bot: Telegraf) => {
    bot.hears(/get\s(\d{2})\/((\d{2})|(\d{1}))/, async (ctx) => {
        const setDay = ctx.match[1] + '/' + ctx.match[2]
        const user = {
            name: ctx.from.username,
            telegramId: ctx.from.id
        }
        const validateUser = await validationUser(user)
        const validateMonth = await validationMonth(validateUser.id, setDay)
        const day = await findOneDay(validateMonth.id, setDay)

        if (day) {
            const formatDate = formatHours(day.workHours)
            ctx.reply(`date: ${day.date} \n time: ${formatDate}
            `)
        }
        else {
            ctx.reply("Нет информации")
        }

    })
}


export const setDayController = async (bot: Telegraf) => {
    bot.hears(/set\s(\d{2})\/((\d{2})|(\d{1}))\s(10|11|12|[1-9]):([0-5][0-9])$/, async (ctx) => {
        const min = convertMinToDecimal(ctx.match[6])
        const setDay = ctx.match[1] + '/' + ctx.match[2]
        const hours = Number(ctx.match[5]) + min
        const user = {
            name: ctx.from.username,
            telegramId: ctx.from.id
        }
        const validateUser = await validationUser(user)
        const validateMonth = await validationMonth(validateUser.id, setDay)
        const newDay = await createDay(validateMonth.id,setDay,hours)
        ctx.reply(JSON.stringify(newDay))
    })
}
