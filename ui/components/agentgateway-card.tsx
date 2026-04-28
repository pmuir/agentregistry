"use client"

import { AgentGateway } from "@/lib/admin-api"
import { Button } from "@/components/ui/button"
import {
  TooltipProvider,
} from "@/components/ui/tooltip"
import { Globe, Trash2 } from "lucide-react"

interface AgentGatewayCardProps {
  gateway: AgentGateway
  onClick?: () => void
  onDelete?: (gateway: AgentGateway) => void
}

function statusColor(status: string): string {
  switch (status) {
    case "healthy":
      return "text-green-600"
    case "unhealthy":
      return "text-red-600"
    default:
      return "text-muted-foreground"
  }
}

function statusDot(status: string): string {
  switch (status) {
    case "healthy":
      return "bg-green-500"
    case "unhealthy":
      return "bg-red-500"
    default:
      return "bg-gray-400"
  }
}

export function AgentGatewayCard({ gateway, onClick, onDelete }: AgentGatewayCardProps) {
  return (
    <TooltipProvider>
      <div
        className="group flex items-start gap-3.5 py-4 px-2 -mx-2 rounded-md cursor-pointer transition-colors hover:bg-muted/50"
        onClick={() => onClick?.()}
      >
        <div className="w-10 h-10 rounded bg-primary/8 flex items-center justify-center flex-shrink-0 mt-0.5">
          <Globe className="h-4 w-4 text-primary" />
        </div>

        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-0.5">
            <h3 className="text-lg font-semibold truncate">{gateway.name}</h3>
            <span className={`flex items-center gap-1.5 text-xs ${statusColor(gateway.status)}`}>
              <span className={`w-2 h-2 rounded-full ${statusDot(gateway.status)}`} />
              {gateway.status}
            </span>
          </div>

          <div className="flex flex-wrap items-center gap-x-3 gap-y-1 text-sm text-muted-foreground">
            <span className="font-mono truncate">{gateway.address}</span>
          </div>
        </div>

        <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0">
          {onDelete && (
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7 text-destructive hover:text-destructive hover:bg-destructive/10"
              onClick={(e) => { e.stopPropagation(); onDelete(gateway) }}
              aria-label="Delete gateway"
            >
              <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
            </Button>
          )}
        </div>
      </div>
    </TooltipProvider>
  )
}
