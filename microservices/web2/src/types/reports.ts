export type Report = {
  id: number | null
  shift: number
  number: number
  receipt: number
  turn: number
  dt: string | null
  success: boolean
  error: string | null
  solution: string | null
  mark: string | null
  responsible: string | null
}