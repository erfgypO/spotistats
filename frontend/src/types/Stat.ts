export enum StatType {
  ARTISTS = "artists",
  // TRACK = "track",
}

export interface Stat {
  name: string;
  percentage: number;
  datapointCount: number;
}
