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
  displayname: string;
  createdAt: string;
  updatedAt: string;
}
