# API Integration Test Script
$baseUrl = "http://localhost:8080"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host " Running Ticket System API Verification   " -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

# 1. Health Check
Write-Host "`n1. Testing GET /health..." -NoNewline
try {
    $health = Invoke-RestMethod -Uri "$baseUrl/health" -Method Get
    if ($health.status -eq "ok") {
        Write-Host " [PASS]" -ForegroundColor Green
    } else {
        Write-Host " [FAIL]" -ForegroundColor Red; exit 1
    }
} catch {
    Write-Host " [FAIL]: $_" -ForegroundColor Red; exit 1
}

# 2. Register Users
Write-Host "`n2. Registering User 1 and User 2..." -NoNewline
$user1Body = @{ email = "alice@example.com"; password = "password123" } | ConvertTo-Json
$user2Body = @{ email = "bob@example.com"; password = "password123" } | ConvertTo-Json

$user1 = Invoke-RestMethod -Uri "$baseUrl/auth/register" -Method Post -Body $user1Body -ContentType "application/json"
$user2 = Invoke-RestMethod -Uri "$baseUrl/auth/register" -Method Post -Body $user2Body -ContentType "application/json"
Write-Host " [PASS]" -ForegroundColor Green

# 3. Login Users
Write-Host "`n3. Logging in User 1 and User 2..." -NoNewline
$login1 = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method Post -Body $user1Body -ContentType "application/json"
$login2 = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method Post -Body $user2Body -ContentType "application/json"
$token1 = $login1.token
$token2 = $login2.token

if ($token1 -and $token2) {
    Write-Host " [PASS]" -ForegroundColor Green
} else {
    Write-Host " [FAIL]: Tokens missing" -ForegroundColor Red; exit 1
}

$headers1 = @{ Authorization = "Bearer $token1" }
$headers2 = @{ Authorization = "Bearer $token2" }

# 4. Create Tickets
Write-Host "`n4. Creating Tickets for Alice..." -NoNewline
$t1Body = @{ title = "Fix Login Bug"; description = "Users cannot login on mobile" } | ConvertTo-Json
$t1 = Invoke-RestMethod -Uri "$baseUrl/tickets" -Method Post -Headers $headers1 -Body $t1Body -ContentType "application/json"
if ($t1.status -eq "open" -and $t1.id) {
    Write-Host " [PASS] Created Ticket ID: $($t1.id)" -ForegroundColor Green
} else {
    Write-Host " [FAIL]" -ForegroundColor Red; exit 1
}

# 5. List Own Tickets
Write-Host "`n5. Listing Alice's Tickets..." -NoNewline
$aliceTickets = Invoke-RestMethod -Uri "$baseUrl/tickets" -Method Get -Headers $headers1
if ($aliceTickets.Count -eq 1 -and $aliceTickets[0].id -eq $t1.id) {
    Write-Host " [PASS] Found 1 ticket owned by Alice" -ForegroundColor Green
} else {
    Write-Host " [FAIL]" -ForegroundColor Red; exit 1
}

Write-Host "`n6. Checking Ownership Isolation (Bob listing tickets)..." -NoNewline
$bobTickets = Invoke-RestMethod -Uri "$baseUrl/tickets" -Method Get -Headers $headers2
if ($bobTickets.Count -eq 0) {
    Write-Host " [PASS] Bob sees 0 tickets" -ForegroundColor Green
} else {
    Write-Host " [FAIL] Bob can see tickets he doesn't own!" -ForegroundColor Red; exit 1
}

# 7. Get Ticket By ID (Forbidden for Bob)
Write-Host "`n7. Bob trying to fetch Alice's ticket directly..." -NoNewline
try {
    Invoke-RestMethod -Uri "$baseUrl/tickets/$($t1.id)" -Method Get -Headers $headers2
    Write-Host " [FAIL] Bob was able to access Alice's ticket!" -ForegroundColor Red; exit 1
} catch {
    if ($_.Exception.Response.StatusCode -eq [System.Net.HttpStatusCode]::NotFound) {
        Write-Host " [PASS] Correctly returned 404 Not Found" -ForegroundColor Green
    } else {
        Write-Host " [FAIL] Unexpected status code: $_" -ForegroundColor Red; exit 1
    }
}

# 8. Status Transition Tests
Write-Host "`n8. Testing status flow: open -> in_progress..." -NoNewline
$patchBody1 = @{ status = "in_progress" } | ConvertTo-Json
$t1_updated = Invoke-RestMethod -Uri "$baseUrl/tickets/$($t1.id)/status" -Method Patch -Headers $headers1 -Body $patchBody1 -ContentType "application/json"
if ($t1_updated.status -eq "in_progress") {
    Write-Host " [PASS]" -ForegroundColor Green
} else {
    Write-Host " [FAIL]" -ForegroundColor Red; exit 1
}

Write-Host "`n9. Testing status flow: in_progress -> closed..." -NoNewline
$patchBody2 = @{ status = "closed" } | ConvertTo-Json
$t1_closed = Invoke-RestMethod -Uri "$baseUrl/tickets/$($t1.id)/status" -Method Patch -Headers $headers1 -Body $patchBody2 -ContentType "application/json"
if ($t1_closed.status -eq "closed") {
    Write-Host " [PASS]" -ForegroundColor Green
} else {
    Write-Host " [FAIL]" -ForegroundColor Red; exit 1
}

Write-Host "`n10. Testing forbidden transition: closed -> open..." -NoNewline
try {
    $patchBody3 = @{ status = "open" } | ConvertTo-Json
    Invoke-RestMethod -Uri "$baseUrl/tickets/$($t1.id)/status" -Method Patch -Headers $headers1 -Body $patchBody3 -ContentType "application/json"
    Write-Host " [FAIL] Reopening closed ticket succeeded!" -ForegroundColor Red; exit 1
} catch {
    if ($_.Exception.Response.StatusCode -eq [System.Net.HttpStatusCode]::BadRequest) {
        Write-Host " [PASS] Correctly rejected with 400 Bad Request" -ForegroundColor Green
    } else {
        Write-Host " [FAIL] Unexpected response: $_" -ForegroundColor Red; exit 1
    }
}

Write-Host "`n==========================================" -ForegroundColor Cyan
Write-Host " ALL API VERIFICATION TESTS PASSED SUCCESSFULLY! " -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Cyan
