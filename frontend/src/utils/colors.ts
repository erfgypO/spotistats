export function rgbToRgba(rgb: string, alpha: number) {
  return rgb.replace(')', `, ${alpha})`).replace('rgb', 'rgba');
}

export const colors = {
  primary: 'rgb(3, 153, 143)',
  secondary: 'rgb(158, 27, 190)',
  tertiary: 'rgb(211, 47, 47)'
}

