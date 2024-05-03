package main

import "net"

// nextIP увеличивает IP-адрес на единицу
func nextIP(ip net.IP) net.IP {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			return ip
		}
	}
	return nil // Переполнение
}

func countIPAddresses(cidr string) (uint64, error) {
	// Разбираем CIDR-нотацию
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return 0, err
	}

	// Получаем маску сети и вычисляем количество адресов по маске
	mask := ipnet.Mask
	ones, _ := mask.Size()
	var count uint64 = 1 << (32 - ones) // 2^(32-бит_маски)

	return count, nil
}

func isValidIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return false
	} else if ip.IsMulticast() {
		return false
	} else if ip.IsPrivate() {
		return false
	}

	return true
}
