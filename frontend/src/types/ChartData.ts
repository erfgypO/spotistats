export interface ChartOptions {
  labels: string[];
  theme?: {
    mode: 'dark' | 'light';
  }
  stroke?: {
    show: boolean;
    width: number;
    colors: string[];
    dashArray: number;
  },
  fill?: {
    opacity: number;
    colors: string[];
  }
}

export interface NamedSeries {
  name: string;
  data: number[];
}
export interface ChartData {
  series: number[] | NamedSeries[];
  chartOptions: ChartOptions;
}
