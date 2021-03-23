import { Day } from "./day.interface";

export interface Month {
    id: string
    userId?: string | null
    fullDate: string 
    name: String
    day?: Day[],
}