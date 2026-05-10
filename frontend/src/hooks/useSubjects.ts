import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import type { UseQueryResult } from '@tanstack/react-query'
import axiosInstance from '../utils/axiosConfig'
import { parseApiError } from '../utils/parseApiError'
import type { Subject, CreateSubjectPayload } from '../types/subjects'

const API_URL = import.meta.env.VITE_API_URL as string + '/api'

export const useGetSubjects = (): UseQueryResult<Subject[], Error> =>
    useQuery<Subject[], Error>({
        queryKey: ['subjects'],
        queryFn: async () => {
            try {
                const { data } = await axiosInstance.get<{ data: Subject[] }>(`${API_URL}/subjects`)
                return data.data ?? []
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        refetchInterval: 5000,
    })

export const useGetSubjectsForTeacher = (teacherId: number): UseQueryResult<Subject[], Error> =>
    useQuery<Subject[], Error>({
        queryKey: ['teacher-subjects', teacherId],
        queryFn: async () => {
            try {
                const { data } = await axiosInstance.get<{ data: Subject[] }>(`${API_URL}/teacher-subject/teacher/${teacherId}/subjects`)
                return data.data ?? []
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        enabled: !!teacherId,
    })

export const useCreateSubject = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (payload: CreateSubjectPayload) => {
            try {
                const { data } = await axiosInstance.post<{ data: Subject }>(`${API_URL}/subject`, payload)
                return data.data
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['subjects'] }),
    })
}

export const useDeleteSubject = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (id: number) => {
            try {
                await axiosInstance.delete(`${API_URL}/subject/${id}`)
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['subjects'] }),
    })
}

