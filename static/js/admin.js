$(document).ready(function()
{
	$('#add_attr').click(function()
	{
		var attr_line = $($('.hidden_attr_field').html());
		attr_line.insertBefore($(this));
		return false;
	});

	$(".shop_grid_small").sortable(
	{
		update: function()
		{
			sort = {};
			$('.shop_grid_small .product_block').each(function()
			{
				var position = $(this).index() + 1;
				var id       = $(this).data('id');

				sort[id] = position;

			});

			$.ajax(
			{
				url: '/admin/update_sort',
				type: 'post',
				data: sort,
				dataType: 'json'
			});
		},
		cancel: ".edit_pick"
	});

	tinymce.init({
      selector: '.tinymce',
      plugins: 'code',
      height: 500
    });
});