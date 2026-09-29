"use client"

import * as React from "react"
import DashTable from "@/components/dashboard/DashTable"
import { appointmentColumns } from "./columns"
import { appointmentData } from "./data"

export default function AppointmentsTable() {
  const [page, setPage] = React.useState(1)
  return (
    <DashTable
      columns={appointmentColumns}
      data={appointmentData}
      getRowId={(row) => row.id}
      pagination={{ page, pageSize: 4, total: 12, itemLabel: "lịch hẹn hôm nay", onPageChange: setPage }}
    />
  )
}
