import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import type { UseQueryResult } from '@tanstack/react-query'
import axiosInstance from '../utils/axiosConfig'
import { parseApiError } from '../utils/parseApiError'
import type { School, CreateSchoolPayload, UpdateSchoolPayload } from '../types/schools'

const API_URL = import.meta.env.VITE_API_URL as string + '/api'

export const useGetSchools = (): UseQueryResult<School[], Error> =>
    useQuery<School[], Error>({
        queryKey: ['schools'],
        queryFn: async () => {
            try {
                const { data } = await axiosInstance.get<{ data: School[] }>(`${API_URL}/schools`)
                return data.data ?? []
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        refetchInterval: 5000,
    })

export const useCreateSchool = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (payload: CreateSchoolPayload) => {
            try {
                const { data } = await axiosInstance.post<{ data: School }>(`${API_URL}/school`, payload)
                return data.data
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['schools'] }),
    })
}

export const useUpdateSchool = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async ({ id, payload }: { id: number; payload: UpdateSchoolPayload }) => {
            try {
                const { data } = await axiosInstance.put<{ data: School }>(`${API_URL}/school/${id}`, payload)
                return data.data
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['schools'] }),
    })
}

export const useDeleteSchool = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (id: number) => {
            try {
                await axiosInstance.delete(`${API_URL}/school/${id}`)
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['schools'] }),
    })
}

