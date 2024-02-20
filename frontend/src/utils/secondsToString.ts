export function secondsToString(seconds: number): string {
  const days = Math.floor(seconds / (3600 * 24));
  const hours = Math.floor(seconds % (3600 * 24) / 3600);
  const minutes = Math.floor(seconds % 3600 / 60);
  const remainingSeconds = Math.floor(seconds % 60);

  let result = '';
  if (days > 0) {
    result += days + 'd ';
  }
  if (hours > 0) {
    result += hours + 'h ';
  }
  if (minutes > 0) {
    result += minutes + 'm ';
  }
  if (remainingSeconds > 0) {
    result += remainingSeconds + 's';
  }
  return result;
}
