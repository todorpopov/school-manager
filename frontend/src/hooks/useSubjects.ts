import { useQuery } from '@tanstack/react-query'
import type { UseQueryResult } from '@tanstack/react-query'
import axiosInstance from '../utils/axiosConfig'
import { parseApiError } from '../utils/parseApiError'
import type { Subject } from '../types/subjects'

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

