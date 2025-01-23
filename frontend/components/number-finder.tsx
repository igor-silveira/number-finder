"use client"

import {useCallback, useState} from "react"
import {Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle} from "@/components/ui/card"
import {Input} from "@/components/ui/input"
import {Button} from "@/components/ui/button"
import {Alert, AlertDescription} from "@/components/ui/alert"
import {Loader2} from "lucide-react"
import {Label} from "@/components/ui/label"
import {Separator} from "@/components/ui/separator"

interface ApiResponse {
    index: number
    number: number
    is_approximate: boolean
}

export default function NumberFinder() {
    const [number, setNumber] = useState("")
    const [thresholdPercentage, setThresholdPercentage] = useState("")
    const [result, setResult] = useState<ApiResponse | null>(null)
    const [error, setError] = useState<string | null>(null)
    const [isLoading, setIsLoading] = useState(false)
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || ""

    const handleSubmit = useCallback(
        async (e: React.FormEvent) => {
            e.preventDefault()
            setIsLoading(true)
            setError(null)
            setResult(null)

            try {
                let url = `${apiUrl}/api/number/${number}`
                if (thresholdPercentage) {
                    url += `?thresholdPercentage=${thresholdPercentage}`
                }

                const response = await fetch(url)
                const data = await response.json()

                if (!response.ok) {
                    if (response.status === 400) {
                        throw new Error(data.message || "Invalid value parameter")
                    } else if (response.status === 404) {
                        throw new Error(data.message || "No value found within acceptable threshold")
                    } else {
                        throw new Error("An unexpected error occurred")
                    }
                }

                setResult(data as ApiResponse)
            } catch (err) {
                setError(err instanceof Error ? err.message : "An unexpected error occurred")
            } finally {
                setIsLoading(false)
            }
        },
        [number, thresholdPercentage, apiUrl],
    )

    return (
        <Card className="w-full max-w-md">
            <CardHeader>
                <CardTitle className="text-2xl font-bold">Number Finder</CardTitle>
                <CardDescription>Find the index of a number in the sequence</CardDescription>
            </CardHeader>
            <CardContent>
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div className="space-y-2">
                        <Label htmlFor="number">Number</Label>
                        <Input
                            id="number"
                            type="number"
                            placeholder="Enter a number"
                            value={number}
                            onChange={(e) => setNumber(e.target.value)}
                            required
                        />
                    </div>
                    <div className="space-y-2">
                        <Label htmlFor="threshold">Threshold Percentage (optional)</Label>
                        <Input
                            id="threshold"
                            type="number"
                            step="0.01"
                            min="0"
                            max="1"
                            placeholder="e.g., 0.1 for 10%"
                            value={thresholdPercentage}
                            onChange={(e) => setThresholdPercentage(e.target.value)}
                        />
                    </div>
                    <Button type="submit" className="w-full" disabled={isLoading}>
                        {isLoading ? (
                            <>
                                <Loader2 className="mr-2 h-4 w-4 animate-spin"/>
                                Searching...
                            </>
                        ) : (
                            "Find"
                        )}
                    </Button>
                </form>
            </CardContent>
            <Separator className="my-4"/>
            <CardFooter>
                <div className="w-full min-h-[120px] flex items-center justify-center">
                    {isLoading ? (
                        <Loader2 className="h-8 w-8 animate-spin text-primary"/>
                    ) : result ? (
                        <div className="space-y-2 w-full">
                            <p className="font-medium">
                                <span className="text-muted-foreground">Index:</span> {result.index}
                            </p>
                            <p className="font-medium">
                                <span className="text-muted-foreground">Value:</span> {result.number}
                            </p>
                            <p className="font-medium">
                                <span
                                    className="text-muted-foreground">Is Approximate:</span> {result.is_approximate ? "Yes" : "No"}
                            </p>
                        </div>
                    ) : error ? (
                        <Alert variant="destructive" className="w-full flex items-center justify-center gap-2">
                            <AlertDescription className="text-center"> {error} </AlertDescription>
                        </Alert>
                    ) : null}
                </div>
            </CardFooter>
        </Card>
    )
}