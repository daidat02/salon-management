"use client"

import * as React from "react"
import { Checkbox } from "@/components/ui/checkbox"
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { cn } from "cn"

export type DashColumn<T> = {
  id: string
  header: string
  headerClassName?: string
  className?: string
  accessorKey?: keyof T
  cell?: (row: T) => React.ReactNode
}

type DashPaginationProps = {
  page: number
  pageSize: number
  total: number
  totalPages?: number
  onPageChange?: (page: number) => void
  itemLabel?: string
}

type DashTableProps<T> = {
  columns: DashColumn<T>[]
  data: T[]
  getRowId?: (row: T, index: number) => string
  selectable?: boolean
  selectedIds?: Set<string>
  onSelectionChange?: (ids: Set<string>) => void
  emptyText?: string
  rowClassName?: (row: T) => string | undefined
  onRowClick?: (row: T) => void
  footer?: React.ReactNode
  pagination?: DashPaginationProps
}

function getPaginationItems(current: number, totalPages: number): (number | "ellipsis")[] {
  if (totalPages <= 7) return Array.from({ length: totalPages }, (_, i) => i + 1)
  if (current <= 3) return [1, 2, 3, "ellipsis", totalPages]
  if (current >= totalPages - 2) return [1, "ellipsis", totalPages - 2, totalPages - 1, totalPages]
  return [1, "ellipsis", current - 1, current, current + 1, "ellipsis", totalPages]
}

export function DashTable<T>({
  columns,
  data,
  getRowId = (_, i) => String(i),
  selectable = true,
  selectedIds,
  onSelectionChange,
  emptyText = "Không có dữ liệu",
  rowClassName,
  onRowClick,
  footer,
  pagination,
}: DashTableProps<T>) {
  const [internalSelected, setInternalSelected] = React.useState<Set<string>>(new Set())
  const isControlled = selectedIds !== undefined && onSelectionChange !== undefined
  const selected = isControlled ? selectedIds! : internalSelected
  const setSelected = isControlled ? onSelectionChange! : setInternalSelected

  const allSelected = data.length > 0 && selected.size === data.length
  const indeterminate = selected.size > 0 && selected.size < data.length

  const toggleAll = (checked: boolean) => {
    if (checked) {
      const all = new Set(data.map((r, i) => getRowId(r, i)))
      setSelected(all)
    } else {
      setSelected(new Set())
    }
  }

  const toggleOne = (id: string, checked: boolean) => {
    const next = new Set(selected)
    if (checked) next.add(id)
    else next.delete(id)
    setSelected(next)
  }

  return (
    <div className="bg-surface overflow-hidden rounded-md border border-outline shadow-sm">
      <Table>
        <TableHeader>
          <TableRow className="bg-surface-container hover:bg-surface-container">
            {selectable && (
              <TableHead className="w-10">
                <Checkbox
                  checked={allSelected}
                  onCheckedChange={(v) => toggleAll(v === true)}
                  aria-label="Chọn tất cả"
                />
              </TableHead>
            )}
            {columns.map((col) => (
              <TableHead key={col.id} className={cn(col.headerClassName, col.className)}>
                {col.header}
              </TableHead>
            ))}
          </TableRow>
        </TableHeader>
        <TableBody>
          {data.length === 0 ? (
            <TableRow>
              <TableCell colSpan={columns.length + (selectable ? 1 : 0)} className="py-10 text-center text-on-surface-variant">
                {emptyText}
              </TableCell>
            </TableRow>
          ) : (
            data.map((row, idx) => {
              const id = getRowId(row, idx)
              const isSelected = selected.has(id)
              return (
                <TableRow
                  key={id}
                  data-state={isSelected ? "selected" : undefined}
                  onClick={() => onRowClick?.(row)}
                  className={cn(onRowClick && "cursor-pointer", rowClassName?.(row))}
                >
                  {selectable && (
                    <TableCell onClick={(e) => e.stopPropagation()}>
                      <Checkbox
                        checked={isSelected}
                        onCheckedChange={(v) => toggleOne(id, v === true)}
                        aria-label={`Chọn ${id}`}
                      />
                    </TableCell>
                  )}
                  {columns.map((col) => {
                    const content = col.cell
                      ? col.cell(row)
                      : col.accessorKey
                        ? (row[col.accessorKey] as unknown as React.ReactNode)
                        : null
                    return (
                      <TableCell key={col.id} className={col.className}>
                        {content}
                      </TableCell>
                    )
                  })}
                </TableRow>
              )
            })
          )}
        </TableBody>
      </Table>
      {pagination ? (
        <div className="bg-surface flex flex-col items-center justify-between gap-3 border-t border-outline px-6 py-3.5 text-xs text-on-surface-variant sm:flex-row">
          <div>
            Hiển thị{" "}
            <span className="font-semibold text-on-surface">
              {(pagination.page - 1) * pagination.pageSize + 1} -{" "}
              {Math.min(pagination.page * pagination.pageSize, pagination.total)}
            </span>{" "}
            trong tổng số{" "}
            <span className="font-semibold text-on-surface">{pagination.total}</span>{" "}
            {pagination.itemLabel ?? ""}
          </div>
          <Pagination className="mx-0 w-auto justify-end">
            <PaginationContent>
              <PaginationItem>
                <PaginationPrevious
                  href="#"
                  onClick={(e) => {
                    e.preventDefault()
                    if (pagination.page > 1) pagination.onPageChange?.(pagination.page - 1)
                  }}
                  aria-disabled={pagination.page <= 1}
                  className={pagination.page <= 1 ? "pointer-events-none opacity-40" : ""}
                  text="Trước"
                />
              </PaginationItem>
              {getPaginationItems(
                pagination.page,
                pagination.totalPages ?? Math.ceil(pagination.total / pagination.pageSize)
              ).map((item, idx) =>
                item === "ellipsis" ? (
                  <PaginationItem key={`e-${idx}`}>
                    <PaginationEllipsis />
                  </PaginationItem>
                ) : (
                  <PaginationItem key={item}>
                    <PaginationLink
                      href="#"
                      isActive={item === pagination.page}
                      onClick={(e) => {
                        e.preventDefault()
                        pagination.onPageChange?.(item)
                      }}
                    >
                      {item}
                    </PaginationLink>
                  </PaginationItem>
                )
              )}
              <PaginationItem>
                <PaginationNext
                  href="#"
                  onClick={(e) => {
                    e.preventDefault()
                    const totalPages = pagination.totalPages ?? Math.ceil(pagination.total / pagination.pageSize)
                    if (pagination.page < totalPages) pagination.onPageChange?.(pagination.page + 1)
                  }}
                  aria-disabled={
                    pagination.page >= (pagination.totalPages ?? Math.ceil(pagination.total / pagination.pageSize))
                  }
                  className={
                    pagination.page >= (pagination.totalPages ?? Math.ceil(pagination.total / pagination.pageSize))
                      ? "pointer-events-none opacity-40"
                      : ""
                  }
                  text="Tiếp theo"
                />
              </PaginationItem>
            </PaginationContent>
          </Pagination>
        </div>
      ) : (
        footer && (
          <div className="bg-surface flex flex-col items-center justify-between gap-3 border-t border-outline px-6 py-3.5 text-xs text-on-surface-variant sm:flex-row">
            {footer}
          </div>
        )
      )}
    </div>
  )
}

export default DashTable
