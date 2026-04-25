import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import type { UseQueryResult } from '@tanstack/react-query'
import axiosInstance from '../utils/axiosConfig'
import { parseApiError } from '../utils/parseApiError'
import type { Principals, CreatePrincipalPayload, UpdatePrincipalPayload } from '../types/principals.ts'

const API_URL = import.meta.env.VITE_API_URL as string + '/api'

export const useGetDirectors = (): UseQueryResult<Principals[], Error> =>
    useQuery<Principals[], Error>({
        queryKey: ['directors'],
        queryFn: async () => {
            try {
                const { data } = await axiosInstance.get<{ data: Principals[] }>(`${API_URL}/directors`)
                return data.data ?? []
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        refetchInterval: 5000,
    })

export const useCreateDirector = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (payload: CreatePrincipalPayload) => {
            try {
                const { data } = await axiosInstance.post<{ data: Principals }>(`${API_URL}/director`, payload)
                return data.data
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['directors'] }),
    })
}

export const useUpdateDirector = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async ({ id, payload }: { id: number; payload: UpdatePrincipalPayload }) => {
            try {
                const { data } = await axiosInstance.put<{ data: Principals }>(`${API_URL}/director/${id}`, payload)
                return data.data
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['directors'] }),
    })
}

export const useDeleteDirector = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (id: number) => {
            try {
                await axiosInstance.delete(`${API_URL}/director/${id}`)
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['directors'] }),
    })
}
