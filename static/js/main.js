var products_map = {};

$(document).ready(function()
{
	$(document).on('click', '.buy_button', function()
	{
		$.fancybox.open($('#order_form_popup'));
		return false;
	});

	$('#order_form').submit(function()
	{
		var data = $(this).serialize();
		$.ajax(
		{
			url: '/create_order',
			data: data,
			method: 'POST',
			success: function()
			{
				$.fancybox.close();
				show_message('Тут будет текст с благодарностью за покупку.');
			},
			error: function()
			{
				show_message('An error occured, please try again later.');
			}
		});
		return false;
	});

	$('.show_mobile_menu').click(function()
	{
		$('.header_nav').slideDown();
		return false;
	});

	$('.close_mobile_menu').click(function()
	{
		$('.header_nav').slideUp();
		return false;
	});

	if ($('.home_page').length)
	{
		// get all products json and save it to products_map
		$.ajax(
		{
			url: '/get_products_json',
			dataType: 'json',
			success: function(json)
			{
				for (i in json['products'])
				{
					var id = json['products'][i]['id'];
					products_map[id] = json['products'][i]
				}
			}
		});

		$('.product_block').click(function(e)
		{
			if (e.target.className != 'buy_button')
			{
				var product_id  = $(this).data('id');
				open_product_popup(product_id);
				return false;
			}
		});

		// image preload:
		$('.product_block').hover(function()
		{
			var product_id  = $(this).data('id');
			$('<img>').attr('src', products_map[product_id]['image']);
		});
	}

	$('.price_value').each(function()
	{
		var price = $(this).text();
		$(this).text(add_spaces(price));
	});
});

function open_product_popup(product_id)
{
	var name        = products_map[product_id]['name'];
	var price       = products_map[product_id]['price'];
	var pick        = products_map[product_id]['image'];
	var description = products_map[product_id]['description'];
	var attributes  = products_map[product_id]['attributes'];
	var url         = products_map[product_id]['url'];

	var attributes_html = '';
	for (key in attributes)
	{
		var attr_name = key;
		var attr_val  = attributes[key];

		attributes_html += `<p>${attr_name}: ${attr_val}</p>`;
	}

	$('#popup_product_name').text(name);
	$('#popup_product_price_value').text(add_spaces(price));
	$('#popup_product_pick').attr("src", "");
	$('#popup_product_pick').attr("src", pick);
	$('#product_popip_pick_link').attr('href', pick);
	$('#popup_product_text').text(description);
	$('#popup_product_attributes').html(attributes_html);

	history.pushState(url, name, url)

	$.fancybox.open(
	{
		src: '#product_popup',
		afterClose: function()
		{
			history.back();
		}
	});
}

function add_spaces(x)
{
    return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, " ");
}

function show_message(message_text)
{
	$.fancybox.open(message_text);
}