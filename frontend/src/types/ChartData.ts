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

export interface RadarChartOptions {
  responsive: boolean;
  maintainAspectRatio: boolean;
  elements?: {
    line?: {
      borderWidth: number;
    }
  },
  plugins?: {
    datalabels?: {
      formatter: (value: any, context: any) => string;
    }
  },
  scales?: {
    r?: {
      grid?: {
        color: string;
      },
      display: boolean;
      ticks?: {
        display: boolean;
      },
      angleLines?: {
        display: boolean;
        color: string;
      },
      suggestedMin: number;
    }
  }
}

export const DefaultRadarChartOptions: RadarChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  elements: {
    line: {
      borderWidth: 3
    }
  },
  scales: {
    r: {
      grid: {
        color: 'rgba(255, 255, 255, 0.2)',
      },
      display: true,
      ticks: {
        display: false,
      },
      angleLines: {
        display: true,
        color: 'rgba(255, 255, 255, 0.2)',
      },
      suggestedMin: 0,
    }
  },
};

export interface NamedSeries {
  name: string;
  data: number[];
}
export interface ChartData {
  series: number[] | NamedSeries[];
  chartOptions: ChartOptions;
}
