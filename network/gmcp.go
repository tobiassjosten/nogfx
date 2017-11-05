package network

/*
EOR = '\025'
ATCP = '\200'
GMCP = '\201'
SE = '\240'
GA = '\249'
SB = '\250'
WILL = '\251'
WONT = '\252'
DO = '\253'
DONT = '\254'
IAC = '\255'

IAC_WILL_ATCP = IAC + WILL + ATCP
IAC_WONT_ATCP = IAC + WONT + ATCP
IAC_DO_ATCP = IAC + DO + ATCP
IAC_DONT_ATCP = IAC + DONT + ATCP
IAC_SB_ATCP = IAC + SB + ATCP
IAC_SE = IAC + SE
IAC_DO_EOR = IAC + DO + EOR
IAC_WILL_EOR = IAC + WILL + EOR
IAC_WONT_EOR = IAC + WONT + EOR
IAC_GA = IAC + GA
*/

/*
p.gmcp('Core.Supports.Set ["Char 1", "Char.Skills 1", "Char.Items 1", "Room 1", "IRE.Composer 1", "IRE.Rift 1"]')
p.gmcp("Char.Skills.Get")
p.gmcp("Char.Items.Inv")
p.gmcp("IRE.Rift.Request")
*/
