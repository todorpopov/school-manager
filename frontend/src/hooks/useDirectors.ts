import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import type { UseQueryResult } from '@tanstack/react-query'
import axiosInstance from '../utils/axiosConfig'
import { parseApiError } from '../utils/parseApiError'
import type { Director, CreateDirectorPayload, UpdateDirectorPayload } from '../types/directors'

const BASE = '/api'

export const useGetDirectors = ():
    UseQueryResult<Director[], Error> => useQuery<Director[], Error>({
        queryKey: ['directors'],
        queryFn: async () => {
            try {
                const { data } = await axiosInstance.get<{ data: Director[] }>(`${BASE}/directors`)
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
        mutationFn: async (payload: CreateDirectorPayload) => {
            try {
                const { data } = await axiosInstance.post<{ data: Director }>(`${BASE}/director`, payload)
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
        mutationFn: async ({ id, payload }: { id: number; payload: UpdateDirectorPayload }) => {
            try {
                const { data } = await axiosInstance.put<{ data: Director }>(`${BASE}/director/${id}`, payload)
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
                await axiosInstance.delete(`${BASE}/director/${id}`)
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['directors'] }),
    })
}


