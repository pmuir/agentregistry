"use client"

import { useState } from "react"
import { AgentGateway } from "@/lib/admin-api"
import { Button } from "@/components/ui/button"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import {
  Globe,
  Code,
  Copy,
  Check,
  Clock,
} from "lucide-react"

interface AgentGatewayDetailProps {
  gateway: AgentGateway
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

export function AgentGatewayDetail({ gateway }: AgentGatewayDetailProps) {
  const [activeTab, setActiveTab] = useState("overview")
  const [jsonCopied, setJsonCopied] = useState(false)

  const handleCopyJson = async () => {
    try {
      await navigator.clipboard.writeText(JSON.stringify(gateway, null, 2))
      setJsonCopied(true)
      setTimeout(() => setJsonCopied(false), 2000)
    } catch (err) {
      console.error('Failed to copy JSON:', err)
    }
  }

  const formatDate = (dateString: string) => {
    try {
      return new Date(dateString).toLocaleString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
      })
    } catch {
      return dateString
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-start gap-4">
        <div className="w-12 h-12 rounded bg-primary/8 flex items-center justify-center flex-shrink-0">
          <Globe className="h-6 w-6 text-primary" />
        </div>
        <div className="flex-1 min-w-0">
          <h1 className="text-2xl font-bold truncate mb-1">{gateway.name}</h1>
          <p className="text-[15px] text-muted-foreground font-mono truncate">{gateway.address}</p>
        </div>
      </div>

      <div className="flex flex-wrap gap-2">
        <span className={`flex items-center gap-1.5 px-2.5 py-1 bg-muted rounded text-sm ${statusColor(gateway.status)}`}>
          <span className={`w-2 h-2 rounded-full ${statusDot(gateway.status)}`} />
          {gateway.status}
        </span>
        {gateway.updatedAt && (
          <span className="flex items-center gap-1.5 px-2.5 py-1 bg-muted rounded text-sm">
            <Clock className="h-3 w-3 text-muted-foreground" />
            {formatDate(gateway.updatedAt)}
          </span>
        )}
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
        <TabsList className="mb-4">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="raw">Raw</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          <section>
            <h3 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground mb-2">Details</h3>
            <div className="space-y-3">
              <div>
                <span className="text-sm text-muted-foreground">ID</span>
                <p className="text-sm font-mono">{gateway.id}</p>
              </div>
              <div>
                <span className="text-sm text-muted-foreground">Address</span>
                <p className="text-sm font-mono">{gateway.address}</p>
              </div>
              <div>
                <span className="text-sm text-muted-foreground">Status</span>
                <p className={`text-sm ${statusColor(gateway.status)}`}>{gateway.status}</p>
              </div>
            </div>
          </section>
        </TabsContent>

        <TabsContent value="raw">
          <div className="rounded-lg border p-4">
            <div className="flex items-center justify-between mb-3">
              <h3 className="text-sm font-semibold flex items-center gap-2">
                <Code className="h-4 w-4" />
                Raw JSON
              </h3>
              <Button variant="outline" size="sm" onClick={handleCopyJson} className="gap-1.5 h-7 text-xs">
                {jsonCopied ? <><Check className="h-3 w-3" /> Copied</> : <><Copy className="h-3 w-3" /> Copy</>}
              </Button>
            </div>
            <pre className="bg-muted p-3 rounded-md overflow-x-auto text-xs leading-relaxed">
              {JSON.stringify(gateway, null, 2)}
            </pre>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
