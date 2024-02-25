export interface Stat {
  name: string;
  percentage: number;
  datapointCount: number;
  spotifyUrl: string;
}

export interface HourlyStat {
  hour: number;
  seconds: number;
  songName: string;
  color: RGB;
}

export interface RGB {
  r: number;
  g: number;
  b: number;
}
