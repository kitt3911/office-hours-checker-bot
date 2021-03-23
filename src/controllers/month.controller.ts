import { isDate } from "moment";
import { Telegraf } from "telegraf";
import { User } from "../interfaces/user.interface";
import { getDayOnMonthService } from "../service/month.service";
import { validationMonth, validationUser } from "../service/validation.service";
import { formatDate, formatHours } from "../utils/formatDate";

export const getThisMonthController = async (bot: Telegraf) => {
    bot.command('show_all_info',async ctx => {
        const userDto:User = {
            name: ctx.from.username,
            telegramId: ctx.from.id
        }
        const user = await validationUser(userDto)
        const month = await validationMonth(user.id)
        const days = await getDayOnMonthService(user.id,month.id)
        let daysToString: string = ""
        days.map(item => {
            daysToString = daysToString + `date: ${item.date} \n time: ${item.workHours ? formatHours(item.workHours) : '0'} \n`
        })
        console.log(days)
        if(days.length > 0){
            ctx.reply(daysToString)
        }
        else{
            ctx.reply("нет информации")
        }
    })
}