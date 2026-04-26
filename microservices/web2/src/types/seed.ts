export type Seed = {
  seed: string
  seed_ru: string
  min_density: number
  max_density: number
  tank_capacity: number
  deleted_at: string
}

export type SeedWithBunker = Seed & {
  amount: number
  bunker: number
}