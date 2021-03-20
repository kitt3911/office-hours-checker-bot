import { PrismaClient } from '@prisma/client'
import { Telegraf } from 'telegraf'
import { botToken } from './config'
import moment from 'moment'
import { formatMonth } from './utils/formatDate'

const prisma = new PrismaClient()

async function main() {
    const bot = new Telegraf(botToken)
    bot.start(async (ctx) => {
        const user = {
            name: ctx.from.username,
            telegramId: ctx.from.id
        }
        let findUser = await prisma.user.findFirst({
            where: { telegramId: user.telegramId }
        })
        if (!findUser) {
            findUser = await prisma.user.create({
                data: { ...user }
            })
        }

        let thisDate = formatMonth()
        const findMonth = await prisma.month.findFirst({
            where: {
                fullDate: thisDate.date,
                user: findUser,
                userId: findUser.id
            }
        })

        if (!findMonth) {
              await prisma.month.create({
                data: {
                    fullDate: thisDate.date,
                    userId: findUser.id,
                    name: thisDate.month
                }
            })
        }

        ctx.reply(`Welcome ${findUser.name}`)
    })

    bot.hears(/(\d{2})\/((\d{2})|(\d{1}))\s\d/, async (ctx) => {
        const setDay = ctx.match[1] + '/' + ctx.match[2]
        const hours = ctx.match[3]
        const formatDate = formatMonth(setDay)
        console.log(formatDate)
        const user = {
            name: ctx.from.username,
            telegramId: ctx.from.id
        }
        let findUser = await prisma.user.findFirst({
            where: { telegramId: user.telegramId }
        })
        if (findUser) {
            let findMonth = await prisma.month.findFirst({
                where: {
                    fullDate: formatDate.date,
                    user: findUser,
                    userId: findUser.id
                }
            })
            if (!findMonth) {
                findMonth = await prisma.month.create({
                    data: {
                        fullDate: formatDate.date,
                        userId: findUser.id,
                        name: formatDate.month
                    }
                })
            }

        }
    })


    bot.launch()
}

main().catch(e => {
    throw e
}).finally(async () => {
    await prisma.$disconnect()
})