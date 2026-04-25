import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import type { UseQueryResult } from '@tanstack/react-query'
import axiosInstance from '../utils/axiosConfig'
import { parseApiError } from '../utils/parseApiError'
import type { Teacher, CreateTeacherPayload, UpdateTeacherPayload } from '../types/teachers'

const API_URL = import.meta.env.VITE_API_URL as string + '/api'

export const useGetTeachers = (): UseQueryResult<Teacher[], Error> =>
    useQuery<Teacher[], Error>({
        queryKey: ['teachers'],
        queryFn: async () => {
            try {
                const { data } = await axiosInstance.get<{ data: Teacher[] }>(`${API_URL}/teachers`)
                return data.data ?? []
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        refetchInterval: 5000,
    })

export const useCreateTeacher = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (payload: CreateTeacherPayload) => {
            try {
                const { data } = await axiosInstance.post<{ data: Teacher }>(`${API_URL}/teacher`, payload)
                return data.data
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['teachers'] }),
    })
}

export const useUpdateTeacher = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async ({ id, payload }: { id: number; payload: UpdateTeacherPayload }) => {
            try {
                const { data } = await axiosInstance.put<{ data: Teacher }>(`${API_URL}/teacher/${id}`, payload)
                return data.data
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['teachers'] }),
    })
}

export const useDeleteTeacher = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (id: number) => {
            try {
                await axiosInstance.delete(`${API_URL}/teacher/${id}`)
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['teachers'] }),
    })
}

