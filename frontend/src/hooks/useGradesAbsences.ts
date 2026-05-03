import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import axiosInstance from '../utils/axiosConfig'
import { parseApiError } from '../utils/parseApiError'
import type { Grade, Absence, CreateGradePayload, CreateAbsencePayload, BulkCreateGradesPayload, BulkCreateAbsencesPayload } from '../types/gradesAbsences'

const API_URL = import.meta.env.VITE_API_URL as string + '/api'

// TODO endpoint
export const useGetGradesForStudent = (studentId: number, enabled = true) =>
    useQuery<Grade[], Error>({
        queryKey: ['student-grades', studentId],
        queryFn: async () => {
            try {
                const { data } = await axiosInstance.get<{ data: Grade[] }>(`${API_URL}/student/${studentId}/grades`)
                return data.data ?? []
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        enabled: enabled && studentId > 0,
    })

// TODO endpoint
export const useGetAbsencesForStudent = (studentId: number, enabled = true) =>
    useQuery<Absence[], Error>({
        queryKey: ['student-absences', studentId],
        queryFn: async () => {
            try {
                const { data } = await axiosInstance.get<{ data: Absence[] }>(`${API_URL}/student/${studentId}/absences`)
                return data.data ?? []
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        enabled: enabled && studentId > 0,
    })

export const useGetAllGrades = () =>
    useQuery<Grade[], Error>({
        queryKey: ['all-grades'],
        queryFn: async () => {
            try {
                const { data } = await axiosInstance.get<{ data: Grade[] }>(`${API_URL}/grades`)
                return data.data ?? []
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        refetchInterval: 5000,
    })

export const useGetAllAbsences = () =>
    useQuery<Absence[], Error>({
        queryKey: ['all-absences'],
        queryFn: async () => {
            try {
                const { data } = await axiosInstance.get<{ data: Absence[] }>(`${API_URL}/absences`)
                return data.data ?? []
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        refetchInterval: 5000,
    })

export const useBulkCreateGrades = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (payload: BulkCreateGradesPayload) => {
            try {
                const { data } = await axiosInstance.post<{ data: Grade[] }>(`${API_URL}/grades/bulk`, payload)
                return data.data
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['all-grades'] }),
    })
}

export const useBulkCreateAbsences = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (payload: BulkCreateAbsencesPayload) => {
            try {
                const { data } = await axiosInstance.post<{ data: Absence[] }>(`${API_URL}/absences/bulk`, payload)
                return data.data
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['all-absences'] }),
    })
}

export const useCreateGrade = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (payload: CreateGradePayload) => {
            try {
                const { data } = await axiosInstance.post<{ data: Grade }>(`${API_URL}/grade`, payload)
                return data.data
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['all-grades'] }),
    })
}

export const useDeleteGrade = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (id: number) => {
            try {
                await axiosInstance.delete(`${API_URL}/grade/${id}`)
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['all-grades'] }),
    })
}

export const useCreateAbsence = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (payload: CreateAbsencePayload) => {
            try {
                const { data } = await axiosInstance.post<{ data: Absence }>(`${API_URL}/absence`, payload)
                return data.data
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['all-absences'] }),
    })
}

export const useDeleteAbsence = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (id: number) => {
            try {
                await axiosInstance.delete(`${API_URL}/absence/${id}`)
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['all-absences'] }),
    })
}
