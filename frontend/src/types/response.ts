export interface LoginResponse {
  accessToken: string;
  expiresAt: number;
}

export interface ErrorResponse {
  error: string;
}

export interface UserResponse {
  id: number;
  username: string;
  displayName: string;
  connectedToSpotify: boolean;
  datapointCount: number;
}
