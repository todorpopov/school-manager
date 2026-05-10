import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import type { UseQueryResult } from '@tanstack/react-query'
import axiosInstance from '../utils/axiosConfig'
import { parseApiError } from '../utils/parseApiError'
import type { Term, CreateTermPayload } from '../types/terms'

const API_URL = import.meta.env.VITE_API_URL as string + '/api'

export const useGetTerms = (): UseQueryResult<Term[], Error> =>
    useQuery<Term[], Error>({
        queryKey: ['terms'],
        queryFn: async () => {
            try {
                const { data } = await axiosInstance.get<{ data: Term[] }>(`${API_URL}/terms`)
                return data.data ?? []
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        refetchInterval: 5000,
    })

export const useCreateTerm = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (payload: CreateTermPayload) => {
            try {
                const { data } = await axiosInstance.post<{ data: Term }>(`${API_URL}/term`, payload)
                return data.data
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['terms'] }),
    })
}

export const useDeleteTerm = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (id: number) => {
            try {
                await axiosInstance.delete(`${API_URL}/term/${id}`)
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['terms'] }),
    })
}

