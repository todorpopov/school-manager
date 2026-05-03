import { useState } from 'react'
import type { ResourceManagerProps } from './types'
import { ResourceForm } from './ResourceForm'
import { ConfirmDelete } from './ConfirmDelete'

type Mode = 'list' | 'create' | 'edit'

export function ResourceManager<T extends { [key: string]: unknown }>({
    title,
    data,
    fields,
    idKey,
    onCreate,
    onUpdate,
    onDelete,
    isLoading = false,
    error = null,
}: ResourceManagerProps<T>) {
    const [mode, setMode] = useState<Mode>('list')
    const [selected, setSelected] = useState<T | null>(null)
    const [deleteTarget, setDeleteTarget] = useState<T | null>(null)

    const tableFields = fields.filter((f) => !f.hideInTable)

    const openCreate  = () => { setSelected(null); setMode('create') }
    const openEdit    = (row: T) => { setSelected(row); setMode('edit') }
    const openDelete  = (row: T) => setDeleteTarget(row)
    const closeDelete = () => setDeleteTarget(null)
    const backToList  = () => { setSelected(null); setMode('list') }

    const handleCreate = async (values: Partial<T>) => { await onCreate(values); backToList() }
    const handleUpdate = async (values: Partial<T>) => { if (!selected) return; await onUpdate(selected[idKey], values); backToList() }
    const handleDelete = async () => { if (!deleteTarget) return; await onDelete(deleteTarget[idKey]); closeDelete() }

    const selectedLabel = selected
        ? (() => {
            const nameField = fields.find((f) => (f.type === 'text' || f.type === 'email') && !f.hideInTable)
            const raw = nameField ? selected[nameField.key] : selected[idKey]
            return typeof raw === 'object' || raw === undefined ? String(selected[idKey]) : String(raw)
        })()
        : ''

    const modeTitle = {
        list:   title,
        create: `New ${title.replace(/s$/, '')}`,
        edit:   `Edit — ${selectedLabel}`,
    }[mode]

    return (
        <div className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg overflow-hidden">
            {deleteTarget && (
                <ConfirmDelete
                    onConfirm={handleDelete}
                    onCancel={closeDelete}
                />
            )}
            <div className="flex items-center justify-between px-5 py-4 bg-slate-50 dark:bg-slate-900 border-b border-slate-200 dark:border-slate-700">
                <h2 className="flex items-center gap-2 m-0 text-lg font-semibold text-slate-800 dark:text-slate-100">
                    {mode !== 'list' && (
                        <button
                            onClick={backToList}
                            className="text-indigo-500 hover:text-indigo-700 text-xl leading-none bg-transparent border-none cursor-pointer p-0"
                        >
                            ←
                        </button>
                    )}
                    {modeTitle}
                </h2>
                {mode === 'list' && (
                    <button
                        onClick={openCreate}
                        className="px-4 py-2 bg-indigo-500 hover:bg-indigo-600 text-white text-sm font-medium rounded-md cursor-pointer border-none transition-colors"
                    >
                        + Add
                    </button>
                )}
            </div>
            {error && (
                <div className="mx-5 mt-3 px-4 py-2 bg-red-100 text-red-700 rounded-md text-sm">
                    {error}
                </div>
            )}
            {mode === 'list' && (
                <div className="overflow-x-auto">
                    {isLoading && (
                        <div className="flex justify-center py-10">
                            <div className="w-7 h-7 border-4 border-slate-200 border-t-indigo-500 rounded-full animate-spin" />
                        </div>
                    )}
                    {!isLoading && data.length === 0 && (
                        <p className="text-center text-slate-400 py-10 text-sm">No records found.</p>
                    )}
                    {!isLoading && data.length > 0 && (
                        <table className="w-full border-collapse text-sm">
                            <thead className="bg-slate-100 dark:bg-slate-900">
                                <tr>
                                    {tableFields.map((f) => (
                                        <th key={String(f.key)} className="px-4 py-3 text-left text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wide whitespace-nowrap">
                                            {f.label}
                                        </th>
                                    ))}
                                    <th className="px-4 py-3 text-right text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wide whitespace-nowrap">
                                        Actions
                                    </th>
                                </tr>
                            </thead>
                            <tbody>
                                {data.map((row) => (
                                    <tr key={String(row[idKey])} className="border-t border-slate-100 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-900/50">
                                        {tableFields.map((f) => (
                                            <td key={String(f.key)} className="px-4 py-3 text-slate-700 dark:text-slate-300 max-w-xs overflow-hidden text-ellipsis whitespace-nowrap">
                                                {f.renderCell
                                                    ? f.renderCell(row[f.key], row)
                                                    : Array.isArray(row[f.key])
                                                        ? (row[f.key] as unknown[]).join(', ')
                                                        : String(row[f.key] ?? '')}
                                            </td>
                                        ))}
                                        <td className="px-4 py-3">
                                            <div className="flex gap-2 justify-end">
                                                <button
                                                    onClick={() => openEdit(row)}
                                                    className="px-3 py-1 text-xs font-medium rounded border border-slate-300 dark:border-slate-600 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 cursor-pointer bg-transparent transition-colors"
                                                >
                                                    Edit
                                                </button>
                                                <button
                                                    onClick={() => openDelete(row)}
                                                    className="px-3 py-1 text-xs font-medium rounded border border-red-200 dark:border-red-800 text-red-500 hover:bg-red-50 dark:hover:bg-red-900/30 cursor-pointer bg-transparent transition-colors"
                                                >
                                                    Delete
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    )}
                </div>
            )}
            {mode === 'create' && (
                <ResourceForm<T>
                    fields={fields}
                    onSubmit={handleCreate}
                    onCancel={backToList}
                    submitLabel="Create"
                    isLoading={isLoading}
                />
            )}
            {mode === 'edit' && selected && (
                <ResourceForm<T>
                    fields={fields}
                    initialValues={selected}
                    onSubmit={handleUpdate}
                    onCancel={backToList}
                    submitLabel="Save changes"
                    isLoading={isLoading}
                    isEdit={true}
                />
            )}
        </div>
    )
}
