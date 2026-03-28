export type FieldType = 'text' | 'email' | 'password' | 'number' | 'select' | 'multiselect' | 'textarea'

export interface SelectOption {
  label: string
  value: string | number
}

export interface FieldConfig<T> {
  /** Key of the field in the data object */
  key: keyof T
  label: string
  type: FieldType
  required?: boolean
  placeholder?: string
  options?: SelectOption[]       // for select / multiselect
  /** Hide field from the table column list */
  hideInTable?: boolean
  /** Hide field from create/edit forms */
  hideInForm?: boolean
  /** Custom render for table cells */
  renderCell?: (value: T[keyof T], row: T) => React.ReactNode
}

export interface ResourceManagerProps<T extends { [key: string]: unknown }> {
  /** Title shown in the header */
  title: string
  /** Array of records to display */
  data: T[]
  /** Field definitions */
  fields: FieldConfig<T>[]
  /** The key used as unique identifier */
  idKey: keyof T
  /** Called with the form values when user submits the create form */
  onCreate: (values: Partial<T>) => Promise<void>
  /** Called with id + form values when user submits the edit form */
  onUpdate: (id: T[keyof T], values: Partial<T>) => Promise<void>
  /** Called with the record id when user confirms deletion */
  onDelete: (id: T[keyof T]) => Promise<void>
  /** Loading state */
  isLoading?: boolean
  /** Optional error message */
  error?: string | null
}

