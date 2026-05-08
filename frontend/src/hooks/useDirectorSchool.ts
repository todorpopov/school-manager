import { useQuery } from '@tanstack/react-query'
import axiosInstance from '../utils/axiosConfig'
import { useAuth } from './useAuth'
import type { Principals } from '../types/principals'

const API_URL = import.meta.env.VITE_API_URL as string + '/api'

export function useDirectorSchoolId(): number | undefined {
    const { user } = useAuth()

    const { data } = useQuery<Principals>({
        queryKey: ['director-by-user', user?.userId],
        queryFn: async () => {
            const { data } = await axiosInstance.get<{ data: Principals }>(`${API_URL}/director/user/${user!.userId}`)
            return data.data
        },
        enabled: !!user?.userId,
    })

    return data?.school?.school_id
}

