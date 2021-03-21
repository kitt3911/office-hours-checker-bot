import { PrismaClient } from '@prisma/client'
import { Telegraf } from 'telegraf'
import { botToken, prisma } from './config'
import { getDayController, setDayController } from './controllers/day.controller'
import { getThisMonthController } from './controllers/month.controller'
import { startController } from './controllers/strat.controller'


async function main() {
    const bot = new Telegraf(botToken)
    await startController(bot)
    await setDayController(bot)
    await getDayController(bot)
    await getThisMonthController(bot)

    bot.launch()
}

main().catch(e => {
    throw e
}).finally(async () => {
    await prisma.$disconnect()
})