import axios from "axios";
import * as Constants from "../../consts/consts.ts";
import { Permissions, UserStatus } from "../../types/authTypes.ts";

export async function getAuthStatus(): Promise<Permissions> {
    try {
        const response = await axios.get<UserStatus>(`${Constants.AUTH_API_URL}/status`);
        console.log('GET: /auth/status');
        console.log('Response data:', response.data);
        return response.data.permissions;
    } catch (error) {
        console.error('Error getting user status', error);
        throw error;
    }
}