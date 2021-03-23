import { PrismaClient } from "@prisma/client"
import { Telegraf } from "telegraf"
import { prisma } from "../config"
import { createDay, findOneDay } from "../service/day.service"
import { validationMonth, validationUser } from "../service/validation.service"
import { formatHours } from "../utils/formatDate"

export const getDayController = async(bot: Telegraf)=>{
    bot.hears(/get\s(\d{2})\/((\d{2})|(\d{1}))/,async (ctx)=> {
        const setDay = ctx.match[1] + '/' + ctx.match[2]
        const user = {
            name: ctx.from.username,
            telegramId: ctx.from.id
        }
        const validateUser = await validationUser(user)
        const validateMonth = await validationMonth(validateUser.id,setDay)
        const day = await findOneDay(validateMonth.id,setDay) 
    
        if(day){
            const formatDate = formatHours(day.workHours)
            ctx.reply(`date: ${day.date} \n time: ${formatDate}
            `)
        }
        else {
            ctx.reply("Нет информации")
        }

    })
}


export const setDayController = async(bot: Telegraf)=>{
    bot.hears(/set\s(\d{2})\/((\d{2})|(\d{1}))\s(10|11|12|[1-9]):([0-5][0-9])$/, async (ctx) => {
        const setDay = ctx.match[1] + '/' + ctx.match[2]
        const hours = ctx.match[5]
        const user = {
            name: ctx.from.username,
            telegramId: ctx.from.id
        }
        const validateUser = await validationUser(user)
        const validateMonth = await validationMonth(validateUser.id,setDay)
       const newDay = await createDay(validateMonth.id,Number(hours),setDay)
       // ctx.reply(JSON.stringify(newDay))
       ctx.reply(ctx.match[5] + ctx.match[6] + ' ')
    })
}