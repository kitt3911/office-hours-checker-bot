import { PrismaClient } from '@prisma/client'
import { Telegraf } from 'telegraf'
import { botToken } from './config'
import { validationMonth, validationUser } from './service/validation.service'

const prisma = new PrismaClient()

async function main() {
    const bot = new Telegraf(botToken)
    bot.start(async (ctx) => {
        const user = {
            name: ctx.from.username,
            telegramId: ctx.from.id
        }
        const validateUser = await validationUser(prisma,user)

        ctx.reply(`Welcome ${validateUser.name}`)
    })

    bot.hears(/(\d{2})\/((\d{2})|(\d{1}))\s\d/, async (ctx) => {
        const setDay = ctx.match[1] + '/' + ctx.match[2]
        const hours = ctx.match[3]
        const user = {
            name: ctx.from.username,
            telegramId: ctx.from.id
        }
        const validateUser = await validationUser(prisma,user)
        const validateMonth = await validationMonth(prisma,setDay,validateUser.id)
        ctx.reply(JSON.stringify(validateMonth))
    })


    bot.launch()
}

main().catch(e => {
    throw e
}).finally(async () => {
    await prisma.$disconnect()
})