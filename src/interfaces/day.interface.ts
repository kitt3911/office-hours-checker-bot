export interface Day {
    monthId?:  string | null
    id?:  string
    date:  string 
    come?: Date | null
    go?: Date | null
    workHours: number
  }