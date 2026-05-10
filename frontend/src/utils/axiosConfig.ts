import axios from 'axios'

const axiosInstance = axios.create()

axiosInstance.interceptors.request.use(
    (config) => {
        config.headers = config.headers ?? {}

        const token = localStorage.getItem('token')
        if (token) {
            config.headers['Authorization'] = token
        }

        const userStr = localStorage.getItem('sm_user')
        if (userStr) {
            try {
                const user = JSON.parse(userStr)
                if (user.sessionId) {
                    config.headers['X-Session-Id'] = user.sessionId
                }
            } catch (e) {
                console.error('Failed to parse user from localStorage:', e)
            }
        }

        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

export default axiosInstance

