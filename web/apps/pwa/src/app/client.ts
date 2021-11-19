import apiClient from "@core/api-client";
import axios from "axios";

export const axiosInstance = axios.create()
export default new apiClient(`${process.env.REACT_APP_SERVER_URL}`, axiosInstance)
